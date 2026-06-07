import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeSnapchat } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { ExtractorFormShell, ExtractorTextField, useExtractorForm } from "@/components/extractor-form";

const snapchatSearchDefaults = {
	owner: "",
	highlight: "",
};

const snapchatSearchSchema = z.object({
	owner: z.string().catch(snapchatSearchDefaults.owner),
	highlight: z.string().catch(snapchatSearchDefaults.highlight),
});

export const Route = createFileRoute("/snapchat")({
	component: Snapchat,
	search: {
		middlewares: [stripSearchParams(snapchatSearchDefaults)],
	},
	validateSearch: snapchatSearchSchema,
});

function Snapchat() {
	const { owner, highlight } = Route.useSearch();
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { owner, highlight },
		validators: {
			onChange: snapchatSearchSchema,
			onSubmit: snapchatSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				highlight: z.string().min(1, "highlight ID is required"),
			}),
		},
		mutation: scrapeSnapchat,
		autoSubmitWhen: ({ owner, highlight }) => owner.length > 0 || highlight.length > 0,
		buildMutationArgs: ({ owner, highlight }) => ({ owner, post: highlight }),
		buildSearch: ({ owner, highlight }) => ({ owner, highlight }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="owner">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="owner"
						placeholder="https://www.snapchat.com/@<OWNER>/highlight/<highlight>"
					/>
				)}
			</form.Field>
			<form.Field name="highlight">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="owner"
						placeholder="https://www.snapchat.com/@<owner>/highlight/<HIGHLIGHT>"
					/>
				)}
			</form.Field>
		</ExtractorFormShell>
	);
}
