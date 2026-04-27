import { createFileRoute } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

export const Route = createFileRoute("/vsco")({
	component: VSCO,
});

function VSCO() {
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
							<FieldLabel>owner</FieldLabel>
							<Input placeholder="https://tiktok.com/@OWNER/video/id" />
						</Field>
						<Field>
							<FieldLabel>post ID</FieldLabel>
							<Input placeholder="https://tiktok.com/@owner/video/ID" />
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
