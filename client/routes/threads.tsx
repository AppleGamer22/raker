import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeThreads } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { ExtractorFormShell, ExtractorTextField, useExtractorForm } from "@/components/extractor-form";

const threadsSearchDefaults = {
	owner: "",
	post: "",
};

const threadsSearchSchema = z.object({
	owner: z.string().catch(threadsSearchDefaults.owner),
	post: z.string().catch(threadsSearchDefaults.post),
});

export const Route = createFileRoute("/threads")({
	component: Threads,
	search: {
		middlewares: [stripSearchParams(threadsSearchDefaults)],
	},
	validateSearch: threadsSearchSchema,
});

function Threads() {
	const { owner, post } = Route.useSearch();
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { owner, post },
		validators: {
			onChange: threadsSearchSchema,
			onSubmit: threadsSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				post: z.string().min(1, "post ID is required"),
			}),
		},
		mutation: scrapeThreads,
		autoSubmitWhen: ({ owner, post }) => owner.length > 0 || post.length > 0,
		buildMutationArgs: ({ owner, post }) => ({ owner, post }),
		buildSearch: ({ owner, post }) => ({ owner, post }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="owner">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="owner"
						placeholder="https://www.threads.com/@OWNER/post/<post>"
					/>
				)}
			</form.Field>
			<form.Field name="post">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="post"
						placeholder="https://www.threads.net/@owner/post/<POST>"
					/>
				)}
			</form.Field>
		</ExtractorFormShell>
	);
}
