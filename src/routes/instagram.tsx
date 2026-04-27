import { createFileRoute } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

export const Route = createFileRoute("/instagram")({
	component: Instagram,
});

function Instagram() {
	return (
		<Card className="m-2 min-h-[calc(100dvh-var(--header-height)-1.1rem)] md:min-h-[calc(100dvh-1rem)]">
			<CardContent>
				<form
					onSubmit={(e) => {
						e.preventDefault();
					}}
				>
					<FieldGroup>
						<Field>
							<FieldLabel>post ID</FieldLabel>
							<Input placeholder="https://www.instagram.com/p/ID" />
						</Field>
					</FieldGroup>
				</form>
			</CardContent>
			<CardFooter>
				<Field orientation="horizontal">
					<Button type="submit" form="bug-report-form">
						Submit
					</Button>
				</Field>
			</CardFooter>
		</Card>
	);
}
