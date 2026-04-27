import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

import { useUser } from "@/hooks/user-provider";

export const Route = createFileRoute("/history")({
	component: History,
});

function History() {
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	return <p>history</p>;
}
