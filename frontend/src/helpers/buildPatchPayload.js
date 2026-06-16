function normalizeForCompare(value) {
	if (value === null || typeof value === "undefined") return "";
	if (typeof value === "string") return value.trim();
	if (typeof value === "object") return JSON.stringify(value);
	return String(value);
}

export function buildPatchPayload(original = {}, current = {}, options = {}) {
	const fields = options.fields ?? Object.keys(current ?? {});
	const transformValue = options.transformValue ?? ((key, value) => value);

	return fields.reduce((payload, key) => {
		const originalValue = normalizeForCompare(original?.[key]);
		const currentValue = normalizeForCompare(current?.[key]);

		if (originalValue === currentValue) {
			return payload;
		}

		payload[key] = transformValue(key, current?.[key], original?.[key]);
		return payload;
	}, {});
}
