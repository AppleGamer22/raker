import { execSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";

import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { defineConfig, type PluginOption } from "vite";

const version = execSync("git describe --tags --abbrev=0").toString().trim();

const MODEL_URL =
	"https://storage.googleapis.com/mediapipe-models/object_detector/efficientdet_lite0/float32/1/efficientdet_lite0.tflite";

export default defineConfig(({ command }) => {
	const isProduction = command === "build";
	return {
		define: {
			"import.meta.env.VITE_GIT_TAG": JSON.stringify(version),
			"import.meta.env.VITE_MODEL_URL": JSON.stringify(`/${path.basename(MODEL_URL)}`),
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
			}) as PluginOption,
			{
				name: "setup-mediapipe-assets",
				async configResolved() {
					const publicDir = path.resolve(import.meta.dirname, "public");
					const wasmDestDir = path.join(publicDir, "wasm");
					const srcDir = path.resolve(import.meta.dirname, "node_modules/@mediapipe/tasks-vision/wasm");
					const modelDest = path.join(publicDir, "efficientdet_lite0.tflite");

					try {
						// 1. Copy WASM files using cp
						execSync(`mkdir -p "${wasmDestDir}" && cp -R "${srcDir}/"* "${wasmDestDir}/"`, {
							stdio: "inherit",
						});
						console.log("Successfully copied WASM files.");

						// 2. Download the model file if it doesn't already exist locally
						if (!fs.existsSync(modelDest)) {
							console.log("Downloading MediaPipe model weights...");
							const response = await fetch(MODEL_URL);

							if (!response.ok) {
								throw new Error(`Failed to download model: ${response.status} ${response.statusText}`);
							}

							const arrayBuffer = await response.arrayBuffer();
							fs.writeFileSync(modelDest, Buffer.from(arrayBuffer));
							console.log("Successfully downloaded model weights to public/");
						} else {
							console.log("Model weights already exist locally, skipping download.");
						}
					} catch (error) {
						console.error("Failed to set up MediaPipe assets:", error);
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
