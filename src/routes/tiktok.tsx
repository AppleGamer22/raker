import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";

const tikTokSearchSchema = z.object({
	owner: z.string().catch(""),
	post: z.string().catch(""),
	incognito: z.boolean().catch(false),
});

export const Route = createFileRoute("/tiktok")({
	component: TikTok,
	validateSearch: tikTokSearchSchema,
});

function TikTok() {
	const { owner, post, incognito } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });

	const form = useForm({
		defaultValues: {
			owner,
			post,
			incognito,
		},
		validators: {
			onChange: tikTokSearchSchema,
			onSubmit: tikTokSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				post: z.string().min(1, "post ID is required"),
			}),
		},
		onSubmit: async ({ value: { owner, post, incognito } }) => {
			await navigate({ search: { owner, post, incognito }, replace: true });
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
										placeholder="https://tiktok.com/@OWNER/video/id"
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
										placeholder="https://tiktok.com/@owner/video/ID"
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
			<CardFooter>{/* TODO: results */}</CardFooter>
		</form>
	);
}
