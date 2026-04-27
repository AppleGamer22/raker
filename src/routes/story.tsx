import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/story")({
	component: Story,
});

function Story() {
	return <p>story</p>;
}
