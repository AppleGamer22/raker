import { SidebarIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/ui/theme-toggle";

export default function Header({ toggleMenu }: { toggleMenu: () => void }) {
	return (
		<header className="sticky top-0 z-50 flex w-full items-center border-b bg-background backdrop-blur-lg">
			<div className="grid h-(--header-height) w-full grid-cols-[1fr_auto_1fr] items-center px-4">
				<div className="justify-self-start">
					<Button className="h-8 w-8" variant="outline" size="icon" onClick={toggleMenu}>
						<SidebarIcon />
					</Button>
				</div>
				<div className="justify-self-center">
					<img alt="Raker Logo" src="/raker.svg" className="w-6" />
				</div>
				<div className="justify-self-end">
					<ThemeToggle />
				</div>
			</div>
		</header>
	);
}
