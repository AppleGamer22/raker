import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/instagram")({
	component: Instagram,
});

function Instagram() {
	return <p>instagram</p>;
}
