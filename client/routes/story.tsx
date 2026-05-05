import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { scrapeStory } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { useUser } from "@/hooks/user-provider";

const storySearchDefaults = {
	owner: "",
};

const storySearchSchema = z.object({
	owner: z.string().catch(storySearchDefaults.owner),
	// incognito: z.boolean().catch(false),
});

export const Route = createFileRoute("/story")({
	component: Story,
	search: {
		middlewares: [stripSearchParams(storySearchDefaults)],
	},
	validateSearch: storySearchSchema,
});

function Story() {
	const { owner } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const storyMutation = useMutation(scrapeStory);
	const [result, setResult] = useState<ScrapeResponse | null>(null);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

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
		onSubmit: async ({ value: { owner } }) => {
			try {
				const result = await storyMutation.mutateAsync({ post: owner });
				setResult(result);
				await navigate({ search: { owner }, replace: true });
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
		if (owner && owner.length > 0) {
			form.handleSubmit();
		}
	}, [form, owner, username]);

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
						<Button type="submit" disabled={storyMutation.isPending} className="mb-3">
							Submit
						</Button>
					</Field>
					{storyMutation.isPending && (
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
