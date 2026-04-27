import { createFileRoute } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";

export const Route = createFileRoute("/tiktok")({
	component: TikTok,
});

function TikTok() {
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
						<Field orientation="horizontal" className="w-fit">
							<FieldLabel>Incognito</FieldLabel>
							<Switch />
						</Field>
					</FieldGroup>
				</form>
			</CardContent>
			<CardFooter>
				<Field orientation="horizontal">
					<Button type="submit">Submit</Button>
				</Field>
				{/* TODO: results */}
			</CardFooter>
		</>
	);
}
