#!/usr/bin/env node
import inquirer from "inquirer";
import { execSync } from "child_process";
import degit from "degit";
import fs from "fs";
import path from "path";
import os from "os";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * CONFIG
 * - Todo está en el MISMO repo:
 *   https://github.com/JulianVNuam/generate-app-cli/tree/main/project-templates
 * - degit NO usa /tree/...; se usa la sintaxis github:user/repo/subdir#branch
 */
const GITHUB_USER = "JulianVNuam";
const GITHUB_REPO = "generate-app-cli";
const SUBDIR_BASE = "project-templates";
const DEFAULT_BRANCH = "main";

// Base para degit (sin /tree/)
const TEMPLATE_REPO_BASE = `github:${GITHUB_USER}/${GITHUB_REPO}/${SUBDIR_BASE}/`;

/** Map de plantillas */
const templateMap = {
  react: { base: "react", tailwind: "react-tailwind" },
  nextjs: { base: "nextjs", tailwind: "nextjs-tailwind" },
};

function detectPkgManager() {
  const ua = process.env.npm_config_user_agent || "";
  if (ua.includes("pnpm")) return "pnpm";
  if (ua.includes("yarn")) return "yarn";
  return "npm";
}

function run(cmd, cwd) {
  execSync(cmd, { stdio: "inherit", cwd });
}

/** Fallback para repos privados o errores de degit */
async function cloneWithGitFallback({ templateName, destDir }) {
  const tmpDir = path.join(os.tmpdir(), `cli-tpl-${Date.now()}`);
  // Clona repo completo vía SSH (requiere tener acceso)
  const repoSsh = `git@github.com:${GITHUB_USER}/${GITHUB_REPO}.git`;
  console.log(
    `\n🔁 Usando fallback: git clone (SSH) → ${repoSsh}#${DEFAULT_BRANCH}`
  );
  run(`git clone --depth=1 -b ${DEFAULT_BRANCH} ${repoSsh} "${tmpDir}"`);

  const from = path.join(tmpDir, SUBDIR_BASE, templateName);
  if (!fs.existsSync(from)) {
    throw new Error(
      `No existe la carpeta de plantilla: ${SUBDIR_BASE}/${templateName} en la rama ${DEFAULT_BRANCH}`
    );
  }
  fs.mkdirSync(destDir, { recursive: true });
  fs.cpSync(from, destDir, { recursive: true });
  fs.rmSync(tmpDir, { recursive: true, force: true });
}

/** Clona plantilla con degit y, si falla, usa git fallback */
async function cloneTemplateOrFallback({ templateName, destDir }) {
  const src = `${TEMPLATE_REPO_BASE}${templateName}#${DEFAULT_BRANCH}`;
  try {
    const emitter = degit(src);
    await emitter.clone(destDir);
  } catch (err) {
    console.warn(
      `\n⚠️  degit falló (${
        err?.code || "UNKNOWN"
      }). Intentando fallback con git clone...`
    );
    await cloneWithGitFallback({ templateName, destDir });
  }
}

async function main() {
  console.log("\n🚀 Bienvenido al Inicializador de Proyectos React/Next.js\n");

  const answers = await inquirer.prompt([
    {
      type: "input",
      name: "projectName",
      message: "Nombre del proyecto:",
      default: "mi-proyecto-app",
    },
    {
      type: "list",
      name: "framework",
      message: "¿Qué framework deseas usar?",
      choices: ["React (Vite)", "Next.js"],
    },
    {
      type: "checkbox",
      name: "features",
      message: "Selecciona las características opcionales:",
      choices: [
        { name: "Tailwind CSS (usa plantilla dedicada)", value: "tailwind" },
        { name: "ESLint + Prettier", value: "eslintPrettier" },
        { name: "Husky + lint-staged", value: "husky" },
        { name: "Jest / Testing Library", value: "testing" },
        { name: "Storybook", value: "storybook" },
        { name: "Internacionalización (i18n)", value: "i18n" },
        { name: "Autenticación (Auth) básica (NextAuth)", value: "auth" },
      ],
    },
  ]);

  const pkgm = detectPkgManager();
  const selectedFramework = answers.framework.includes("React")
    ? "react"
    : "nextjs";
  const wantTailwind = answers.features.includes("tailwind");
  const templateName = wantTailwind
    ? templateMap[selectedFramework].tailwind
    : templateMap[selectedFramework].base;

  const targetDir = path.join(process.cwd(), answers.projectName);

  console.log(`\n📦 Descargando plantilla: ${templateName}...`);
  await cloneTemplateOrFallback({ templateName, destDir: targetDir });

  // Guarda config elegida
  fs.writeFileSync(
    path.join(targetDir, "config.cli.json"),
    JSON.stringify(answers, null, 2)
  );

  console.log("\n📥 Instalando dependencias...");
  run(`${pkgm} install`, targetDir);

  // ===== Features opcionales =====
  if (answers.features.includes("eslintPrettier")) {
    console.log("\n🧹 Configurando ESLint + Prettier...");
    run(
      `${pkgm} ${
        pkgm === "yarn" ? "add -D" : "install -D"
      } eslint prettier eslint-config-prettier eslint-plugin-react eslint-plugin-react-hooks`,
      targetDir
    );
    const isNext = selectedFramework === "nextjs";
    const eslintrc = {
      extends: [
        "eslint:recommended",
        isNext ? "next/core-web-vitals" : "plugin:react/recommended",
        "prettier",
      ],
      plugins: ["react", "react-hooks"],
      parserOptions: { ecmaVersion: "latest", sourceType: "module" },
      settings: isNext ? {} : { react: { version: "detect" } },
    };
    fs.writeFileSync(
      path.join(targetDir, ".eslintrc.json"),
      JSON.stringify(eslintrc, null, 2)
    );
    fs.writeFileSync(
      path.join(targetDir, ".prettierrc"),
      JSON.stringify({ semi: true, singleQuote: false }, null, 2)
    );
    fs.writeFileSync(
      path.join(targetDir, ".prettierignore"),
      `node_modules\n.next\ndist\nbuild\n`
    );
  }

  if (answers.features.includes("husky")) {
    console.log("\n🐶 Configurando Husky + lint-staged...");
    run(
      `${pkgm} ${pkgm === "yarn" ? "add -D" : "install -D"} husky lint-staged`,
      targetDir
    );
    run(`npx husky init`, targetDir);
    const hookPath = path.join(targetDir, ".husky", "pre-commit");
    const hook = `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"
${pkgm} exec lint-staged
`;
    fs.writeFileSync(hookPath, hook, { mode: 0o755 });

    const pkgJsonPath = path.join(targetDir, "package.json");
    const pkg = JSON.parse(fs.readFileSync(pkgJsonPath, "utf-8"));
    pkg["lint-staged"] = {
      "*.{js,jsx,ts,tsx}": ["eslint --fix", "prettier --write"],
    };
    fs.writeFileSync(pkgJsonPath, JSON.stringify(pkg, null, 2));
  }

  if (answers.features.includes("testing")) {
    console.log("\n🧪 Agregando Testing...");
    const isReact = selectedFramework === "react";
    if (isReact) {
      run(
        `${pkgm} ${
          pkgm === "yarn" ? "add -D" : "install -D"
        } vitest @testing-library/react @testing-library/jest-dom jsdom`,
        targetDir
      );
      const vitestCfg = `import { defineConfig } from 'vitest/config'
export default defineConfig({ test: { environment: 'jsdom' } })
`;
      fs.writeFileSync(path.join(targetDir, "vitest.config.js"), vitestCfg);
    } else {
      run(
        `${pkgm} ${
          pkgm === "yarn" ? "add -D" : "install -D"
        } @testing-library/react @testing-library/jest-dom jest jest-environment-jsdom ts-jest`,
        targetDir
      );
      const jestCfg = `/** @type {import('jest').Config} */
const config = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
};
module.exports = config;
`;
      fs.writeFileSync(path.join(targetDir, "jest.config.cjs"), jestCfg);
      fs.writeFileSync(
        path.join(targetDir, "jest.setup.js"),
        `import '@testing-library/jest-dom'`
      );
    }
  }

  if (answers.features.includes("storybook")) {
    console.log("\n📚 Instalando Storybook...");
    run(`npx storybook@latest init`, targetDir);
  }

  if (answers.features.includes("i18n")) {
    console.log("\n🌍 Agregando i18n...");
    const isNext = selectedFramework === "nextjs";
    if (isNext) {
      run(
        `${pkgm} ${
          pkgm === "yarn" ? "add" : "install"
        } next-i18next react-i18next i18next`,
        targetDir
      );
    } else {
      run(
        `${pkgm} ${pkgm === "yarn" ? "add" : "install"} react-i18next i18next`,
        targetDir
      );
    }
  }

  if (answers.features.includes("auth")) {
    console.log("\n🔐 Agregando autenticación (NextAuth)...");
    run(`${pkgm} ${pkgm === "yarn" ? "add" : "install"} next-auth`, targetDir);
  }

  console.log("\n✅ Proyecto creado exitosamente en:", answers.projectName);
  console.log("\n👉 Siguiente paso:");
  console.log(`cd ${answers.projectName} && ${pkgm} run dev`);
}

main().catch((e) => {
  console.error("\n❌ Error inesperado:", e?.message || e);
  process.exit(1);
});
