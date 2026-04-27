import { ComputerIcon, MoonIcon, SunIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme, type Theme } from "@/hooks/theme-provider";

function ThemeIcon({ theme }: { theme: Theme }) {
	switch (theme) {
		case "light":
			return <SunIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />;
		case "dark":
			return <MoonIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />;
		case "system":
			return <ComputerIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />;
	}
}

export function ThemeToggle() {
	const { theme, setTheme } = useTheme();

	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="outline" size="icon">
					<ThemeIcon theme={theme} />
					<span className="sr-only">Toggle theme</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end">
				<DropdownMenuItem onClick={() => setTheme("light")}>
					<SunIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
					Light
				</DropdownMenuItem>
				<DropdownMenuItem onClick={() => setTheme("dark")}>
					<MoonIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
					Dark
				</DropdownMenuItem>
				<DropdownMenuItem onClick={() => setTheme("system")}>
					<ComputerIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
					System
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
