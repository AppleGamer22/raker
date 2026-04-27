import { createFileRoute } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";

export const Route = createFileRoute("/highlight")({
	component: Highlight,
});

function Highlight() {
	return (
		<>
			<CardContent>
				<form
					onSubmit={(e) => {
						e.preventDefault();
					}}
				>
					<FieldGroup>
						<Field>
							<FieldLabel>highlight ID</FieldLabel>
							<Input placeholder="https://www.instagram.com/stories/find/highlights/ID" />
						</Field>
						<Field orientation="horizontal" className="w-fit">
							<FieldLabel>Incognito</FieldLabel>
							<Switch />
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
				{/* TODO: results */}
			</CardFooter>
		</>
	);
}
