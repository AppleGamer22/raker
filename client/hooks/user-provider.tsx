import { createContext, useContext, useState, type ReactNode } from "react";

type UserProviderState = {
	username: string | null;
	categories: string[];
};

const UserProviderContext = createContext<UserProviderState | undefined>(undefined);

function readCookie(name: string): string | null {
	if (typeof document === "undefined") {
		return null;
	}

	for (const cookie of document.cookie.split(";")) {
		const trimmedCookie = cookie.trim();
		if (!trimmedCookie) {
			continue;
		}

		const equalsIndex = trimmedCookie.indexOf("=");
		if (equalsIndex === -1) {
			continue;
		}

		const cookieName = trimmedCookie.slice(0, equalsIndex);
		if (cookieName !== name) {
			continue;
		}

		return trimmedCookie.slice(equalsIndex + 1);
	}

	return null;
}

function decodeBase64Url(encodedValue: string): string | null {
	const normalizedValue = encodedValue.replaceAll("-", "+").replaceAll("_", "/");
	const paddedValue = normalizedValue.padEnd(Math.ceil(normalizedValue.length / 4) * 4, "=");

	try {
		return atob(paddedValue);
	} catch {
		return null;
	}
}

function readUsernameFromJwtCookie(): string | null {
	const jwt = readCookie("jwt");
	if (!jwt) {
		return null;
	}

	const tokenParts = jwt.split(".");
	if (tokenParts.length < 2) {
		return null;
	}

	const payloadJson = decodeBase64Url(tokenParts[1]);
	if (!payloadJson) {
		return null;
	}

	try {
		const payload = JSON.parse(payloadJson) as { username?: unknown };
		return typeof payload.username === "string" && payload.username ? payload.username : null;
	} catch {
		return null;
	}
}

export function UserProvider({ children }: { children: ReactNode }) {
	const [username] = useState<string | null>(() => readUsernameFromJwtCookie());
	return <UserProviderContext.Provider value={{ username, categories: [] }}>{children}</UserProviderContext.Provider>;
}

export function useUser() {
	const context = useContext(UserProviderContext);

	if (context === undefined) {
		throw new Error("useUser must be used within a UserProvider");
	}

	return context;
}
