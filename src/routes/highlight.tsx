import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/highlight")({
	component: Highlight,
});

function Highlight() {
	return <p>highlight</p>;
}
