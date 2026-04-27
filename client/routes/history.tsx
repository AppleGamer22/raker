import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/history")({
	component: History,
});

function History() {
	return <p>history</p>;
}
