import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

const highlightSearchSchema = z.object({
	highlight: z.string().catch(""),
});

export const Route = createFileRoute("/highlight")({
	component: Highlight,
	validateSearch: highlightSearchSchema,
});

function Highlight() {
	const { highlight } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });

	const form = useForm({
		defaultValues: {
			highlight,
		},
		validators: {
			onChange: highlightSearchSchema,
			onSubmit: highlightSearchSchema.extend({
				highlight: z.string().min(1, "highlight ID is required"),
			}),
		},
		onSubmit: async ({ value: { highlight } }) => {
			await navigate({ search: { highlight }, replace: true });
		},
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
					<form.Field name="highlight">
						{(field) => {
							const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
							return (
								<Field>
									<FieldLabel htmlFor={field.name}>highlight ID</FieldLabel>
									<Input
										name={field.name}
										value={field.state.value}
										onBlur={field.handleBlur}
										aria-invalid={isInvalid}
										onChange={(e) => field.handleChange(e.target.value)}
										placeholder="https://www.instagram.com/stories/find/highlights/ID"
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
