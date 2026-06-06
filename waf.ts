// @ts-ignore
import { DOMParser } from "https://deno.land/x/deno_dom/deno-dom-wasm.ts";

// @ts-ignore
const html = Deno.args[1];

const document = new DOMParser().parseFromString(html, "text/html");

// Add a simple cookie jar to capture document.cookie writes in the sandbox.
const cookieJar: Map<string, string> = new Map<string, string>();
Object.defineProperty(document, "cookie", {
	get: () => [...cookieJar.entries()].map(([key, value]) => `${key}=${value}`).join("; "),
	set: (raw: string) => {
		const firstPart = raw.split(";")[0] ?? "";
		const eqIndex = firstPart.indexOf("=");
		if (eqIndex === -1) {
			return;
		}

		const name = firstPart.slice(0, eqIndex).trim();
		const value = firstPart.slice(eqIndex + 1).trim();
		if (!name) {
			return;
		}
		if (!value) {
			return;
		}
		cookieJar.set(name, value);
	},
	configurable: true,
});

// Set the page URL
// @ts-ignore
const url = new URL(Deno.args[0]);
Object.defineProperty(document, "URL", {
	value: url.href,
	configurable: true,
});
Object.defineProperty(document, "documentURI", {
	value: url.href,
	configurable: true,
});

// Create a minimal window/location stub for scripts that expect a browser.
const location = {
	href: url.href,
	assign: (next: string) => {
		location.href = next;
	},
	replace: (next: string) => {
		location.href = next;
	},
	reload: () => {},
};

const windowStub = {
	...globalThis,
	document,
	location,
	navigator: globalThis.navigator,
};

Object.defineProperty(document, "defaultView", {
	value: windowStub,
	configurable: true,
});
Object.defineProperty(globalThis, "location", {
	value: location,
	configurable: true,
});

// Create eval context with necessary globals
const scope = {
	window: windowStub,
	document,
	navigator: globalThis.navigator,
	location,
};

async function waitForCookieWrites(timeoutMs = 30e3) {
	if (cookieJar.size > 0) {
		return;
	}

	const startedAt = Date.now();
	while (Date.now() - startedAt < timeoutMs) {
		await new Promise((resolve) => setTimeout(resolve, 10));
		if (cookieJar.size > 0) {
			return;
		}
	}
}

// Execute all <script> tags in order (inline and external).
// Skip non-JS types (e.g. application/json configs).
const scripts = Array.from(document.getElementsByTagName("script")) as HTMLOrSVGScriptElement[];

for (const scriptEl of scripts) {
	const type = scriptEl.getAttribute("type") || "";
	if (type && type !== "text/javascript" && type !== "application/javascript") {
		continue;
	}

	// Expose currentScript while executing so scripts that rely on it work.
	try {
		Object.defineProperty(document, "currentScript", { value: scriptEl, configurable: true });
	} catch {
		(document as any).currentScript = scriptEl;
	}

	let code = "";
	const src = scriptEl.getAttribute("src") || scriptEl.getAttribute("data-src");
	if (src) {
		try {
			const srcUrl = new URL(src, url.href).href;
			code = await fetch(srcUrl).then((r) => r.text());
		} catch {
			// If fetching fails, skip this script but keep going.
			try {
				delete (document as any).currentScript;
			} catch {}
			continue;
		}
	} else {
		code = scriptEl.textContent || "";
	}

	if (!code) {
		try {
			delete (document as any).currentScript;
		} catch {}
		continue;
	}

	try {
		const maybePromise = new Function(...Object.keys(scope), code)(...Object.values(scope));
		if (maybePromise && typeof (maybePromise as PromiseLike<unknown>).then === "function") {
			await maybePromise;
		}
		delete (document as any).currentScript;
	} catch {}
}

await waitForCookieWrites();

console.log(
	JSON.stringify(
		[...cookieJar.entries()].filter(([, value]) => value.length > 0).map(([name, value]) => ({ name, value })),
		null,
		2,
	),
);
