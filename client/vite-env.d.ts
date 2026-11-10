/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_GIT_TAG: string;
	readonly VITE_MODEL_URL: string;
	readonly DEV: boolean;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
