import {
  HttpSuccessResponse,
  HttpSuccessResponseMultiple,
} from "../types/types";
import config from "./config";

export async function getLatestTasks(): Promise<HttpSuccessResponseMultiple> {
  const failedResponse = {
    success: false,
    tasks: [],
  } as HttpSuccessResponseMultiple;
  const url = new URL(config.API_BASE_URL + config.ENDPOINT_SUFFIX.LATEST);
  url.searchParams.append("count", config.NUM_OF_RECENT_TASKS);

  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error("Network response is not ok");
    }
    return response.json();
  } catch (error) {
    return Promise.resolve(failedResponse);
  }
}

export async function createTask(
  formData: FormData,
): Promise<HttpSuccessResponse> {
  const failedResponse = {
    success: false,
    task: [],
  } as HttpSuccessResponse;
  const url = new URL(config.API_BASE_URL + config.ENDPOINT_SUFFIX.CREATE);
  const fetchOptions = {
    method: "POST",
    body: formData,
  };

  try {
    const response = await fetch(url, fetchOptions);
    if (!response.ok) {
      throw new Error("Network response is not ok");
    }
    return response.json();
  } catch (error) {
    return Promise.resolve(failedResponse);
  }
}
