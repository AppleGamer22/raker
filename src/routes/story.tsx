import { useForm } from "@tanstack/react-form";
import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

const storySearchSchema = z.object({
	owner: z.string().catch(""),
	// incognito: z.boolean().catch(false),
});

export const Route = createFileRoute("/story")({
	component: Story,
	validateSearch: storySearchSchema,
});

function Story() {
	const { owner } = Route.useSearch();

	const form = useForm({
		defaultValues: {
			owner,
		},
		validators: {
			onChange: storySearchSchema,
			onSubmit: storySearchSchema.extend({
				owner: z.string().min(1, "story owner is required"),
			}),
		},
		onSubmit: async () => {},
	});

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
				form.handleSubmit();
			}}
		>
			<CardContent>
				<FieldGroup>
					<form.Field name="owner">
						{(field) => {
							const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
							return (
								<Field>
									<FieldLabel htmlFor={field.name}>owner</FieldLabel>
									<Input
										name={field.name}
										value={field.state.value}
										onBlur={field.handleBlur}
										aria-invalid={isInvalid}
										onChange={(e) => field.handleChange(e.target.value)}
										placeholder="https://www.instagram.com/stories/OWNER"
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							);
						}}
					</form.Field>
					<Field orientation="horizontal">
						<Button type="submit">Submit</Button>
					</Field>
				</FieldGroup>
			</CardContent>
			<CardFooter>{/* TODO: results */}</CardFooter>
		</form>
	);
}
