import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState, useRef } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeSnapchat } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
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
	const [result, setResult] = useState<ScrapeResponse | null>(null);

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
				const result = await snapchatMutation.mutateAsync({ owner, post: highlight });
				setResult(result);
				await navigate({ search: { owner, highlight }, replace: true });
			} catch (err) {
				toast.error((err as Error).message, {
					position: "top-center",
				});
			}
		},
	});

	// submit once on initial page load if search params are present
	const initialSubmit = useRef(true);
	useEffect(() => {
		if (!initialSubmit.current) return;
		initialSubmit.current = false;
		if (username === null) return;
		if ((owner && owner.length > 0) || (highlight && highlight.length > 0)) {
			form.handleSubmit();
		}
	}, [form, owner, highlight, username]);

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
						<Button type="submit" className="mb-3">
							Submit
						</Button>
					</Field>
				</FieldGroup>
			</CardContent>
			{result && (
				<CardFooter>
					<Result result={result} setResult={setResult} />
				</CardFooter>
			)}
		</form>
	);
}
