import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { CardContent } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { useUser } from "@/hooks/user-provider";

const instagramSearchSchema = z.object({
	post: z.string().catch(""),
	incognito: z.boolean().catch(false),
});

export const Route = createFileRoute("/instagram")({
	component: Instagram,
	validateSearch: instagramSearchSchema,
});

function Instagram() {
	const { post, incognito } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	const form = useForm({
		defaultValues: {
			post,
			incognito,
		},
		validators: {
			onChange: instagramSearchSchema,
			onSubmit: instagramSearchSchema.extend({
				post: z.string().min(1, "post ID is required"),
			}),
		},
		onSubmit: async ({ value: { post, incognito } }) => {
			await navigate({ search: { post, incognito }, replace: true });
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
										placeholder="https://www.instagram.com/p/ID"
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							);
						}}
					</form.Field>
					<form.Field name="incognito">
						{(field) => (
							<Field orientation="horizontal" className="w-fit">
								<FieldLabel htmlFor={field.name}>Incognito</FieldLabel>
								<Switch
									name={field.name}
									checked={field.state.value}
									onCheckedChange={field.handleChange}
								/>
							</Field>
						)}
					</form.Field>
					<Field orientation="horizontal">
						<Button type="submit">Submit</Button>
					</Field>
				</FieldGroup>
			</CardContent>
			{/* TODO: results */}
			{/* <CardFooter></CardFooter> */}
		</form>
	);
}
