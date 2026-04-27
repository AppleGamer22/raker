import { createFileRoute } from "@tanstack/react-router";
// import { createRootRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: AuthPage, ssr: false });

function AuthPage() {
	return <p>auth</p>;
}
