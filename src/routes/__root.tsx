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
		<>
			{isMobile ? <Header toggleMenu={() => setOpen(true)} /> : null}
			{isMobile ? (
				<MobileMenu open={open} onOpenChange={setOpen} />
			) : (
				<Menu>
					<Outlet />
				</Menu>
			)}
			{isMobile ? <Outlet /> : null}
		</>
	);
}
