import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeInstagram } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { Switch } from "@/components/ui/switch";
import { useUser } from "@/hooks/user-provider";

const instagramSearchDefaults = {
	post: "",
	incognito: false,
};

const instagramSearchSchema = z.object({
	post: z.string().catch(instagramSearchDefaults.post),
	incognito: z.boolean().catch(instagramSearchDefaults.incognito),
});

export const Route = createFileRoute("/instagram")({
	component: Instagram,
	search: {
		middlewares: [stripSearchParams(instagramSearchDefaults)],
	},
	validateSearch: instagramSearchSchema,
});

function Instagram() {
	const { post, incognito } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const instagramMutation = useMutation(scrapeInstagram);
	const [result, setResult] = useState<ScrapeResponse | null>(null);

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
			try {
				const result = await instagramMutation.mutateAsync({ post, incognito });
				setResult(result);
				await navigate({ search: { post, incognito }, replace: true });
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
		if (post && post.length > 0) {
			form.handleSubmit();
		}
	}, [form, post, username]);

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
										id={field.name}
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
									id={field.name}
									name={field.name}
									checked={field.state.value}
									onCheckedChange={field.handleChange}
								/>
							</Field>
						)}
					</form.Field>
					<Field orientation="horizontal">
						<Button type="submit" disabled={instagramMutation.isPending} className="mb-3 w-full sm:w-auto">
							Submit
						</Button>
					</Field>
					{instagramMutation.isPending && (
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
