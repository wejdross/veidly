import { getSupportedLanguage } from "../locale";

import {properGetInstructor } from "./instructor.api"

const { API_URL } = require("../conf");
const { gettoken } = require("../helpers");
const { xhr } = require("./api");

export function getInstructor() {
    return properGetInstructor()
}
