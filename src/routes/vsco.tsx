import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/vsco")({
	component: VSCO,
});

function VSCO() {
	return <p>vsco</p>;
}
