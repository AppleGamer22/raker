import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeStory } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { ExtractorFormShell, ExtractorTextField, useExtractorForm } from "@/components/extractor-form";

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
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { owner },
		validators: {
			onChange: storySearchSchema,
			onSubmit: storySearchSchema.extend({
				owner: z.string().min(1, "story owner is required"),
			}),
		},
		mutation: scrapeStory,
		autoSubmitWhen: ({ owner }) => owner.length > 0,
		buildMutationArgs: ({ owner }) => ({ post: owner }),
		buildSearch: (_, result) => ({ owner: result.post }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="owner">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="owner"
						placeholder="https://www.instagram.com/stories/OWNER"
					/>
				)}
			</form.Field>
		</ExtractorFormShell>
	);
}
