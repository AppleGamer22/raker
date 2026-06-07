import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeHighlight } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { ExtractorFormShell, ExtractorTextField, useExtractorForm } from "@/components/extractor-form";

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
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { highlight },
		validators: {
			onChange: highlightSearchSchema,
			onSubmit: highlightSearchSchema.extend({
				highlight: z.string().min(1, "highlight ID is required"),
			}),
		},
		mutation: scrapeHighlight,
		autoSubmitWhen: ({ highlight }) => highlight.length > 0,
		buildMutationArgs: ({ highlight }) => ({ post: highlight }),
		buildSearch: ({ highlight }) => ({ highlight }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="highlight">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="highlight ID"
						placeholder="https://www.instagram.com/stories/find/highlights/ID"
					/>
				)}
			</form.Field>
		</ExtractorFormShell>
	);
}
