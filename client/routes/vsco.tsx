import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeVSCO } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { ExtractorFormShell, ExtractorTextField, useExtractorForm } from "@/components/extractor-form";

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
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { owner, post },
		validators: {
			onChange: vscoSearchSchema,
			onSubmit: vscoSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				post: z.string().min(1, "post ID is required"),
			}),
		},
		mutation: scrapeVSCO,
		autoSubmitWhen: ({ owner, post }) => owner.length > 0 || post.length > 0,
		buildMutationArgs: ({ owner, post }) => ({ owner, post }),
		buildSearch: ({ owner, post }) => ({ owner, post }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="owner">
				{(field) => (
					<ExtractorTextField field={field} label="owner" placeholder="https://vsco.co/OWNER/media/id" />
				)}
			</form.Field>
			<form.Field name="post">
				{(field) => (
					<ExtractorTextField field={field} label="post ID" placeholder="https://vsco.co/USERNAME/media/ID" />
				)}
			</form.Field>
		</ExtractorFormShell>
	);
}
