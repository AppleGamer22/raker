import { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
	children: React.ReactNode;
	defaultTheme?: Theme;
	storageKey?: string;
};

type ThemeProviderState = {
	theme: Theme;
	computedTheme: Exclude<Theme, "system">;
	setTheme: (theme: Theme) => void;
};

const initialState: ThemeProviderState = {
	theme: "system",
	computedTheme: "light",
	setTheme: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
	children,
	defaultTheme = "system",
	storageKey = "vite-ui-theme",
	...props
}: ThemeProviderProps) {
	const getComputedTheme = (): Exclude<Theme, "system"> => {
		if (typeof window === "undefined") return "light";
		return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
	};

	const [theme, setTheme] = useState<Theme>(() => {
		const storedTheme = localStorage.getItem(storageKey);
		if (storedTheme === "dark" || storedTheme === "light" || storedTheme === "system") {
			return storedTheme;
		}

		return defaultTheme;
	});
	const [computedTheme, setComputedTheme] = useState<Exclude<Theme, "system">>(getComputedTheme);

	useEffect(() => {
		const root = window.document.documentElement;
		const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

		const applyTheme = () => {
			root.classList.remove("light", "dark");

			if (theme === "system") {
				const systemTheme = mediaQuery.matches ? "dark" : "light";
				root.classList.add(systemTheme);
				setComputedTheme(systemTheme);
				return;
			}

			root.classList.add(theme);
			setComputedTheme(theme);
		};

		applyTheme();

		if (theme !== "system") {
			return;
		}

		mediaQuery.addEventListener("change", applyTheme);

		return () => {
			mediaQuery.removeEventListener("change", applyTheme);
		};
	}, [theme]);

	const value = {
		theme,
		computedTheme,
		setTheme: (theme: Theme) => {
			if (typeof window !== "undefined") {
				window.localStorage.setItem(storageKey, theme);
			}
			setTheme(theme);
		},
	};

	return (
		<ThemeProviderContext.Provider {...props} value={value}>
			{children}
		</ThemeProviderContext.Provider>
	);
}

export const useTheme = () => {
	const context = useContext(ThemeProviderContext);

	if (context === undefined) throw new Error("useTheme must be used within a ThemeProvider");

	return context;
};
