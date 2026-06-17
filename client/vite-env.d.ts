/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_GIT_TAG: string;
	readonly DEV: boolean;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
