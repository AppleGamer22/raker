import { execSync } from "node:child_process";
import path from "node:path";

import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const version = execSync("git describe --tags --abbrev=0").toString().trim();

export default defineConfig(({ command }) => {
	const isProduction = command === "build";
	return {
		define: {
			"import.meta.env.VITE_GIT_TAG": JSON.stringify(version),
		},
		resolve: {
			tsconfigPaths: true,
		},
		plugins: [
			tailwindcss(),
			tanstackRouter({
				target: "react",
				autoCodeSplitting: isProduction,
				routesDirectory: "./client/routes",
				generatedRouteTree: "./client/routeTree.gen.ts",
			}),
			devtools(),
			react(),
			reactCompilerPreset(),
			babel({
				presets: [],
			}),
			{
				name: "copy-mediapipe-wasm-cp",
				configResolved() {
					// Destination directory inside public so Vite serves it out of the box
					const destDir = path.resolve(import.meta.dirname, "public/wasm");
					const srcDir = path.relative(import.meta.dirname, "node_modules/@mediapipe/tasks-vision/wasm");

					try {
						execSync(`mkdir -p "${destDir}" && cp -R "${srcDir}/"* "${destDir}/"`, { stdio: "inherit" });
					} catch (error) {
						console.error("[vite-plugin] Failed to copy WASM files:", error);
					}
				},
			},
		],
		server: {
			hmr: {
				overlay: true,
			},
			proxy: {
				"/api": {
					target: "http://localhost:4100",
					changeOrigin: true,
					// rewrite: (path) => path.replace(/^\/api/, ""),
				},
			},
			host: "0.0.0.0",
			allowedHosts: [".ts.net"],
		},
		build: {
			sourcemap: true,
			outDir: "vdist",
		},
		css: {
			devSourcemap: true,
		},
	};
});
