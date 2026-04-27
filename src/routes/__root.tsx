import { createRootRoute, Outlet } from "@tanstack/react-router";
import { useState } from "react";

import Header from "@/components/header";
import { Menu, MobileMenu } from "@/components/menu";
import { useIsMobile } from "@/hooks/use-mobile";

export const Route = createRootRoute({
	component: Root,
	ssr: false,
});

function Root() {
	const [open, setOpen] = useState(false);
	const isMobile = useIsMobile();

	return (
		<main className="w-full max-w-full overflow-x-hidden">
			{isMobile ? <Header toggleMenu={() => setOpen(true)} /> : null}
			{isMobile ? (
				<MobileMenu open={open} onOpenChange={setOpen} />
			) : (
				<Menu>
					<Outlet />
				</Menu>
			)}
			{isMobile ? (
				<div className="w-full max-w-full overflow-x-hidden">
					<Outlet />
				</div>
			) : null}
		</main>
	);
}
