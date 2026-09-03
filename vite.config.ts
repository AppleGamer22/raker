import { execSync } from "node:child_process";

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
			babel({
				presets: [reactCompilerPreset()],
			}),
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
