import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/tiktok")({
	component: TikTok,
});

function TikTok() {
	return <p>tiktok</p>;
}
