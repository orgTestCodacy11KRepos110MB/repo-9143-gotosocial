/*
	GoToSocial
	Copyright (C) 2021-2022 GoToSocial Authors admin@gotosocial.org

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

"use strict";

const syncpipe = require("syncpipe");

module.exports = function useFormSubmit(form, [mutationQuery, result], { changedOnly = true } = {}) {
	return [
		function submitForm(e) {
			e.preventDefault();

			// transform the field definitions into an object with just their values 
			let updatedFields = [];
			const mutationData = syncpipe(form, [
				(_) => Object.values(_),
				(_) => _.map((field) => {
					if (!changedOnly || field.hasChanged()) {
						updatedFields.push(field);
						return [field.name, field.value];
					} else {
						return null;
					}
				}),
				(_) => _.filter((value) => value != null),
				(_) => Object.fromEntries(_)
			]);

			return mutationQuery(mutationData);
		},
		result
	];
};