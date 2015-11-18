#include "include/capi/cef_task_capi.h"

void CEF_CALLBACK cef_task_t_execute(struct _cef_task_t* self) {
    go_TaskExecute(self);
}

void initialize_task(struct _cef_task_t* self) {
    go_Log("initialize_task");
    self->base.size = sizeof(cef_task_t);
    initialize_cef_base((cef_base_t*) self, "task");
    // callbacks
    self->execute = cef_task_t_execute;
}

void helper_cef_post_task(cef_thread_id_t threadId, struct _cef_task_t* what) {
    cef_post_task(threadId, what);
}