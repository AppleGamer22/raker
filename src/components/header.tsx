import { ModeToggle } from "./ui/theme-toggle";

export default function Header() {
	return (
		<header className="sticky top-0 z-50 flex w-full items-center border-b bg-background backdrop-blur-lg">
			<div className="flex h-(--header-height) w-full items-center gap-2 px-4">
				<ModeToggle />
			</div>
		</header>
	);
}
