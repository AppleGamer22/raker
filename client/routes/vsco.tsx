import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeVSCO } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { useUser } from "@/hooks/user-provider";

const vscoSearchDefaults = {
	owner: "",
	post: "",
};

const vscoSearchSchema = z.object({
	owner: z.string().catch(vscoSearchDefaults.owner),
	post: z.string().catch(vscoSearchDefaults.post),
});

export const Route = createFileRoute("/vsco")({
	component: VSCO,
	search: {
		middlewares: [stripSearchParams(vscoSearchDefaults)],
	},
	validateSearch: vscoSearchSchema,
});

function VSCO() {
	const { owner, post } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const vscoMutation = useMutation(scrapeVSCO);
	const [result, setResult] = useState<ScrapeResponse | null>(null);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

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
		onSubmit: async ({ value: { owner, post } }) => {
			try {
				const result = await vscoMutation.mutateAsync({ owner, post });
				setResult(result);
				await navigate({ search: { owner, post }, replace: true });
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
		if ((owner && owner.length > 0) || (post && post.length > 0)) {
			form.handleSubmit();
		}
	}, [form, owner, post, username]);

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
										id={field.name}
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
										id={field.name}
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
						<Button type="submit" disabled={vscoMutation.isPending} className="mb-3 w-full sm:w-auto">
							Submit
						</Button>
					</Field>
					{vscoMutation.isPending && (
						<Field>
							<Progress value={null} className="pb-2" />
						</Field>
					)}
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
