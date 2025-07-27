import typescriptEslint from '@typescript-eslint/eslint-plugin';
import prettier from 'eslint-plugin-prettier';
import parser from 'vue-eslint-parser';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import js from '@eslint/js';
import {FlatCompat} from '@eslint/eslintrc';
import globals from 'globals';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
	baseDirectory: __dirname,
	recommendedConfig: js.configs.recommended,
	allConfig: js.configs.all,
});

export default [
	{
		ignores: ['**/node_modules', '**/dist', 'src/**/*.js', '**/*.tsbuildinfo'],
	},
	...compat.extends(
		'plugin:vue/strongly-recommended',
		'eslint:recommended',
	),
	{
		plugins: {
			'@typescript-eslint': typescriptEslint,
			'prettier': prettier,
		},

		languageOptions: {
			parser,
			ecmaVersion: 2021, // Modern ECMAScript version
			sourceType: 'module', // Use "module" for ES6+ support
			parserOptions: {
				parser: '@typescript-eslint/parser',
				extraFileExtensions: ['.vue'], // Ensure Vue files are parsed
			},
			globals: {
				...globals.browser,
			},
		},

		rules: {
			'vue/no-multiple-template-root': 'off', // Only relevant for Vue 2
			'vue/multi-word-component-names': 'off',
			'@typescript-eslint/no-unused-vars': 'warn', // Better as a warning
			'prettier/prettier': 'off', // Enforce Prettier formatting
			'@typescript-eslint/explicit-module-boundary-types': 'warn', // Encourage explicit types
			'@typescript-eslint/no-explicit-any': 'warn', // Discourage 'any' usage
		},
	},
];
