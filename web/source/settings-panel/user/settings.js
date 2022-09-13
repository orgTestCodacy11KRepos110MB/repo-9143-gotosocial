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

const Promise = require("bluebird");
const React = require("react");
const Redux = require("react-redux");

const api = require("../lib/api");
const user = require("../redux/reducers/user").actions;

const Languages = require("../components/languages");
const Submit = require("../components/submit");

const {
	Checkbox,
	Select,
} = require("../components/form-fields").formFields(user.setProfileVal, (state) => state.user.profile);

module.exports = function UserSettings() {
	const dispatch = Redux.useDispatch();

	const [errorMsg, setError] = React.useState("");
	const [statusMsg, setStatus] = React.useState("");

	function submit() {
		setStatus("PATCHing");
		setError("");
		return Promise.try(() => {
			return dispatch(api.user.updateSettings());
		}).then(() => {
			setStatus("Saved!");
		}).catch((e) => {
			setError(e.message);
			setStatus("");
		});
	}

	return (
		<div className="user-settings">
			<h1>Post settings</h1>
			<Select id="language" name="Default post language">
				<Languages/>
			</Select>
			<Select id="privacy" name="Default post privacy" options={
				<>
					<option value="private">Private / followers-only)</option>
					<option value="unlisted">Unlisted</option>
					<option value="public">Public</option>
				</>
			}>
				<a href="https://docs.gotosocial.org/en/latest/user_guide/posts/#privacy-settings" target="_blank" className="moreinfolink" rel="noreferrer">Learn more about post privacy settings (opens in a new tab)</a>
			</Select>
			<Select id="format" name="Default post format" options={
				<>
					<option value="plain">Plain (default)</option>
					<option value="markdown">Markdown</option>
				</>
			}>
				<a href="https://docs.gotosocial.org/en/latest/user_guide/posts/#input-types" target="_blank" className="moreinfolink" rel="noreferrer">Learn more about post format settings (opens in a new tab)</a>
			</Select>
			<Checkbox
				id="sensitive"
				name="Mark my posts as sensitive by default"
			/>

			<Submit onClick={submit} label="Save post settings" errorMsg={errorMsg} statusMsg={statusMsg}/>
			<hr/>
			<PasswordChange/>
		</div>
	);
};

function PasswordChange() {
	const dispatch = Redux.useDispatch();

	const [errorMsg, setError] = React.useState("");
	const [statusMsg, setStatus] = React.useState("");

	const [oldPassword, setOldPassword] = React.useState("");
	const [newPassword, setNewPassword] = React.useState("");
	const [newPasswordConfirm, setNewPasswordConfirm] = React.useState("");

	function submit() {
		if (newPassword !== newPasswordConfirm) {
			setError("New password and confirm new password did not match!");
			return;
		}
		
		setStatus("PATCHing");
		setError("");
		return Promise.try(() => {
			let data = {
				old_password: oldPassword,
				new_password: newPassword
			};
			return dispatch(api.apiCall("POST", "/api/v1/user/password_change", data, "form"));
		}).then(() => {
			setStatus("Saved!");
			setOldPassword("");
			setNewPassword("");
			setNewPasswordConfirm("");
		}).catch((e) => {
			setError(e.message);
			setStatus("");
		});
	}

	return (
		<>
			<h1>Change password</h1>
			<div className="labelinput">
				<label htmlFor="password">Current password</label>
				<input name="password" id="password" type="password" autoComplete="current-password" value={oldPassword} onChange={(e) => setOldPassword(e.target.value)} />
			</div>
			<div className="labelinput">
				<label htmlFor="new-password">New password</label>
				<input name="new-password" id="new-password" type="password" autoComplete="new-password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} />
			</div>
			<div className="labelinput">
				<label htmlFor="confirm-new-password">Confirm new password</label>
				<input name="confirm-new-password" id="confirm-new-password" type="password" autoComplete="new-password" value={newPasswordConfirm} onChange={(e) => setNewPasswordConfirm(e.target.value)} />
			</div>
			<Submit onClick={submit} label="Save new password" errorMsg={errorMsg} statusMsg={statusMsg}/>
		</>
	);
}