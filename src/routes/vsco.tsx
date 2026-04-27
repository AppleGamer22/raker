import { useForm } from "@tanstack/react-form";
import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

const vscoSearchSchema = z.object({
	owner: z.string().catch(""),
	post: z.string().catch(""),
});

export const Route = createFileRoute("/vsco")({
	component: VSCO,
	validateSearch: vscoSearchSchema,
});

function VSCO() {
	const { owner, post } = Route.useSearch();

	const form = useForm({
		defaultValues: {
			owner,
			post,
		},
		validators: {
			onChange: vscoSearchSchema,
			onSubmit: vscoSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				post: z.string().min(1, "post ID is required"),
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
										placeholder="https://vsco.co/OWNER/media/id"
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							);
						}}
					</form.Field>
					<form.Field name="post">
						{(field) => {
							const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
							return (
								<Field>
									<FieldLabel htmlFor={field.name}>post ID</FieldLabel>
									<Input
										name={field.name}
										value={field.state.value}
										onBlur={field.handleBlur}
										aria-invalid={isInvalid}
										onChange={(e) => field.handleChange(e.target.value)}
										placeholder="https://vsco.co/USERNAME/media/ID"
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
