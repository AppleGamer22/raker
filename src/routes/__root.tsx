import { createRootRoute, Outlet } from "@tanstack/react-router";

import Header from "@/components/header";
import { Menu } from "@/components/menu";
import { useIsMobile } from "@/hooks/use-mobile";

export const Route = createRootRoute({
	component: Root,
	ssr: false,
});

function Root() {
	const isMobile = useIsMobile();

	return (
		<>
			{isMobile ? <Header /> : null}
			<Menu />
			<Outlet />
		</>
	);
}
