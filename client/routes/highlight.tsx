import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeHighlight } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { useUser } from "@/hooks/user-provider";

const highlightSearchDefaults = {
	highlight: "",
};

const highlightSearchSchema = z.object({
	highlight: z.string().catch(highlightSearchDefaults.highlight),
});

export const Route = createFileRoute("/highlight")({
	component: Highlight,
	search: {
		middlewares: [stripSearchParams(highlightSearchDefaults)],
	},
	validateSearch: highlightSearchSchema,
});

function Highlight() {
	const { highlight } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const highlightMutation = useMutation(scrapeHighlight);
	const [result, setResult] = useState<ScrapeResponse | null>(null);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

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
			try {
				const result = await highlightMutation.mutateAsync({ post: highlight });
				setResult(result);
				await navigate({ search: { highlight }, replace: true });
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
		if (highlight && highlight.length > 0) {
			form.handleSubmit();
		}
	}, [form, highlight, username]);

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
										id={field.name}
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
						<Button type="submit" disabled={highlightMutation.isPending} className="mb-3">
							Submit
						</Button>
					</Field>
					{highlightMutation.isPending && (
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
