import { ComputerIcon, MoonIcon, SunIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/hooks/theme-provider";

function ThemeIcon() {
	const { theme, computedTheme } = useTheme();
	switch (theme) {
		case "light":
			return <SunIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />;
		case "dark":
			return <MoonIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />;
		case "system":
			return (
				<>
					<ComputerIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />(
					{computedTheme === "light" ? (
						<SunIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
					) : (
						<MoonIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
					)}
					)
				</>
			);
	}
}

export function ThemeToggle() {
	const { setTheme } = useTheme();

	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="outline">
					<ThemeIcon />
					<span className="sr-only">Toggle theme</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="mt-1" align="end">
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
