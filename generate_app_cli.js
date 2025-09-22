import inquirer from "inquirer";
import { execSync } from "child_process";
import degit from "degit";
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const TEMPLATE_REPO_BASE =
  "github:JulianVNuam/generate-app-cli/project-templates/";

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
  const emitter = degit(`${TEMPLATE_REPO_BASE}${templateName}`);
  await emitter.clone(targetDir);

  fs.writeFileSync(
    path.join(targetDir, "config.cli.json"),
    JSON.stringify(answers, null, 2)
  );

  console.log("\n📥 Instalando dependencias...");
  run(`${pkgm} install`, targetDir);

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
    // activar husky
    run(`npx husky init`, targetDir);
    // agregar hook pre-commit
    const hookPath = path.join(targetDir, ".husky", "pre-commit");
    const hook = `#!/bin/sh\n. "$(dirname "$0")/_/husky.sh"\n${pkgm} exec lint-staged\n`;
    fs.writeFileSync(hookPath, hook, { mode: 0o755 });

    // agregar config lint-staged al package.json
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
      const vitestCfg = `import { defineConfig } from 'vitest/config'\nexport default defineConfig({ test: { environment: 'jsdom' } })\n`;
      fs.writeFileSync(path.join(targetDir, "vitest.config.js"), vitestCfg);
    } else {
      run(
        `${pkgm} ${
          pkgm === "yarn" ? "add -D" : "install -D"
        } @testing-library/react @testing-library/jest-dom jest jest-environment-jsdom ts-jest`,
        targetDir
      );
      const jestCfg = `/** @type {import('jest').Config} */\nconst config = {\n  testEnvironment: 'jsdom',\n  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],\n};\nmodule.exports = config;\n`;
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

main();
