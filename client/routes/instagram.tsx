import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeInstagram } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import {
	ExtractorFormShell,
	ExtractorSwitchField,
	ExtractorTextField,
	useExtractorForm,
} from "@/components/extractor-form";

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
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { post, incognito },
		validators: {
			onChange: instagramSearchSchema,
			onSubmit: instagramSearchSchema.extend({
				post: z.string().min(1, "post ID is required"),
			}),
		},
		mutation: scrapeInstagram,
		autoSubmitWhen: ({ post }) => post.length > 0,
		buildMutationArgs: ({ post, incognito }) => ({ post, incognito }),
		buildSearch: ({ post, incognito }) => ({ post, incognito }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="post">
				{(field) => (
					<ExtractorTextField field={field} label="post ID" placeholder="https://www.instagram.com/p/ID" />
				)}
			</form.Field>
			<form.Field name="incognito">
				{(field) => <ExtractorSwitchField field={field} label="Incognito" />}
			</form.Field>
		</ExtractorFormShell>
	);
}
