# 🚀 Company App Starter CLI

Un **CLI personalizado** para inicializar proyectos en **React (Vite)** o **Next.js**, con características opcionales como **Tailwind CSS**, **ESLint + Prettier**, **Husky + lint-staged**, **Testing**, **Storybook**, **i18n** y **NextAuth**.

La idea es estandarizar y acelerar la creación de nuevos proyectos en la empresa, con configuraciones y dependencias listas desde el inicio.

---

## ✨ Características

- 🔧 Selección de framework: **React (Vite)** o **Next.js**
- 🎨 Opción de incluir **Tailwind CSS** desde plantillas preconfiguradas
- 🧹 Configuración rápida de **ESLint + Prettier**
- 🐶 Hooks de Git con **Husky + lint-staged**
- 🧪 Testing con **Vitest / Testing Library** (React) o **Jest** (Next.js)
- 📚 **Storybook** integrado para documentación de UI
- 🌍 Soporte para **Internacionalización (i18n)**
- 🔐 Autenticación con **NextAuth** (solo Next.js)

---

## 📦 Instalación

Clona el repositorio y enlázalo globalmente:

```bash
git clone https://github.com/JulianVNuam/generate-app-cli.git
cd generate-app-cli
npm install
npm link
```

Esto registrará el comando `company-cli` en tu sistema.

---

## 🚀 Uso

Ejecuta:

```bash
company-cli
```

El CLI te guiará con preguntas:

1. Nombre del proyecto
2. Framework a usar (**React** o **Next.js**)
3. Características opcionales (Tailwind, ESLint, Testing, etc.)

Ejemplo:

```bash
? Nombre del proyecto: demo-next-tailwind
? ¿Qué framework deseas usar? Next.js
? Selecciona las características opcionales: (x) Tailwind CSS, (x) ESLint + Prettier, (x) Husky + lint-staged
```

Esto generará un proyecto **Next.js + Tailwind** listo para ejecutar.

---

## 📂 Estructura de plantillas

El CLI clona plantillas desde el repositorio:

```
project-templates/
├── react/           # React con Vite
├── nextjs/          # Next.js App Router
├── react-tailwind/  # React con Tailwind preconfigurado
└── nextjs-tailwind/ # Next.js con Tailwind preconfigurado
```

Según tus selecciones, se descarga la plantilla correcta y se añaden las configuraciones extra.

---

## 🛠 Desarrollo

Si quieres modificar el CLI:

1. Edita `bin/index.js` (punto de entrada).
2. Reinstala el link global:
   ```bash
   npm unlink -g generate-app-cli
   npm link
   ```
3. Prueba de nuevo con `company-cli`.

---

## 📋 Roadmap

- [ ] Agregar banderas no interactivas (`--framework nextjs --features tailwind,eslintPrettier`)
- [ ] Incluir soporte para monorepo con Turborepo
- [ ] Agregar CI/CD base (GitHub Actions)
- [ ] Soporte Docker en plantillas

---

## 📄 Licencia

MIT © 2025 – Creado por **Julian Vivas / nuam**
