import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeSnapchat } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { Button } from "@/components/ui/button";
import { CardContent } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { useUser } from "@/hooks/user-provider";

const snapchatSearchSchema = z.object({
	owner: z.string().catch(""),
	highlight: z.string().catch(""),
});

export const Route = createFileRoute("/snapchat")({
	component: Snapchat,
	validateSearch: snapchatSearchSchema,
});

function Snapchat() {
	const { owner, highlight } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const snapchatMutation = useMutation(scrapeSnapchat);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	const form = useForm({
		defaultValues: {
			owner,
			highlight,
		},
		validators: {
			onChange: snapchatSearchSchema,
			onSubmit: snapchatSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				highlight: z.string().min(1, "highlight ID is required"),
			}),
		},
		onSubmit: async ({ value: { owner, highlight } }) => {
			try {
				await navigate({ search: { owner, highlight }, replace: true });
				const result = await snapchatMutation.mutateAsync({ owner, post: highlight });
				console.log(result);
			} catch (err) {
				toast.error((err as Error).message, {
					position: "top-center",
				});
			}
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
										placeholder="https://www.snapchat.com/@<OWNER>/highlight/<highlight>"
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							);
						}}
					</form.Field>
					<form.Field name="highlight">
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
										placeholder="https://www.snapchat.com/@<owner>/highlight/<HIGHLIGHT>"
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
			{/* TODO: results */}
			{/* <CardFooter></CardFooter> */}
		</form>
	);
}
