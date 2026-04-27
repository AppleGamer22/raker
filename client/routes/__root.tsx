import { createRootRoute, Outlet } from "@tanstack/react-router";
import { useState } from "react";

import Header from "@/components/header";
import { Menu, MobileMenu } from "@/components/menu";
import { Card } from "@/components/ui/card";
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
			<MobileMenu open={open && isMobile} onOpenChange={setOpen} />

			<Menu>
				<Card className="m-2 min-h-[calc(100dvh-var(--header-height)-1.1rem)] md:min-h-[calc(100dvh-1rem)]">
					<Outlet />
				</Card>
			</Menu>
		</main>
	);
}
