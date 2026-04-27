import { ComputerIcon, MoonIcon, SunIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuLabel,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme, type Theme } from "@/hooks/theme-provider";

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
	const { setTheme, theme } = useTheme();

	return (
		<DropdownMenu>
			<DropdownMenuTrigger render={<Button variant="outline" />}>
				<ThemeIcon />
				<span className="sr-only">Toggle theme</span>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="mt-1" align="end">
				<DropdownMenuGroup>
					<DropdownMenuLabel>Theme</DropdownMenuLabel>
					<DropdownMenuRadioGroup value={theme} onValueChange={(value) => setTheme(value as Theme)}>
						<DropdownMenuRadioItem value="light">
							<SunIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
							Light
						</DropdownMenuRadioItem>
						<DropdownMenuRadioItem value="dark">
							<MoonIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
							Dark
						</DropdownMenuRadioItem>
						<DropdownMenuRadioItem value="system">
							<ComputerIcon className="h-[1.2rem] w-[1.2rem] scale-100 transition-all" />
							System
						</DropdownMenuRadioItem>
					</DropdownMenuRadioGroup>
				</DropdownMenuGroup>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
