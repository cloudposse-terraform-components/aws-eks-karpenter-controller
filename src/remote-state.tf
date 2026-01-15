locals {
  # Remote state is enabled when var.remote_state_enabled is true and the module is enabled
  remote_state_enabled = local.enabled && var.remote_state_enabled

  # Validation: when remote_state_enabled is false, all direct EKS cluster variables must be provided
  _validate_direct_vars = !var.remote_state_enabled ? (
    var.eks_cluster_id != null &&
    var.eks_cluster_arn != null &&
    var.eks_cluster_endpoint != null &&
    var.eks_cluster_certificate_authority_data != null &&
    var.eks_cluster_identity_oidc_issuer != null &&
    var.karpenter_node_role_arn != null
  ) : true
}

# Validation resource to ensure proper configuration
resource "terraform_data" "validate_eks_config" {
  count = local.enabled ? 1 : 0

  lifecycle {
    precondition {
      condition     = local._validate_direct_vars
      error_message = <<-EOT
      When remote_state_enabled is false, all direct EKS cluster variables must be provided:
        - eks_cluster_id
        - eks_cluster_arn
        - eks_cluster_endpoint
        - eks_cluster_certificate_authority_data
        - eks_cluster_identity_oidc_issuer
        - karpenter_node_role_arn
      EOT
    }
  }
}

module "eks" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  # When bypass is true, skip remote state lookup and return defaults
  # We bypass when remote_state_enabled is false (using direct variables instead)
  bypass = !local.remote_state_enabled

  component = var.eks_component_name

  context = module.this.context

  # When bypassed, use direct variables as defaults
  # When not bypassed but remote state is missing, fall back to "deleted" for graceful degradation
  defaults = {
    eks_cluster_id                         = coalesce(var.eks_cluster_id, "deleted")
    eks_cluster_arn                        = coalesce(var.eks_cluster_arn, "deleted")
    eks_cluster_endpoint                   = coalesce(var.eks_cluster_endpoint, "deleted")
    eks_cluster_certificate_authority_data = coalesce(var.eks_cluster_certificate_authority_data, "")
    eks_cluster_identity_oidc_issuer       = coalesce(var.eks_cluster_identity_oidc_issuer, "deleted")
    karpenter_iam_role_arn                 = coalesce(var.karpenter_node_role_arn, "deleted")
  }
}
